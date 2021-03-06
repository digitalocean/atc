package dbng

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/concourse/atc"
)

//go:generate counterfeiter . Resource

type Resource interface {
	ID() int
	Name() string
	PipelineName() string
	Type() string
	Source() atc.Source
	CheckEvery() string
	Tags() atc.Tags
	CheckError() error
	Paused() bool
	WebhookToken() string
	FailingToCheck() bool

	Pause() error
	Unpause() error

	Reload() (bool, error)
}

var resourcesQuery = psql.Select("r.id, r.name, r.config, r.check_error, r.paused, r.pipeline_id, p.name, r.nonce").
	From("resources r").
	Join("pipelines p ON p.id = r.pipeline_id").
	Where(sq.Eq{"r.active": true})

type resource struct {
	id           int
	name         string
	pipelineID   int
	pipelineName string
	type_        string
	source       atc.Source
	checkEvery   string
	tags         atc.Tags
	checkError   error
	paused       bool
	webhookToken string

	conn Conn
}

type ResourceNotFoundError struct {
	Name string
}

func (e ResourceNotFoundError) Error() string {
	return fmt.Sprintf("resource '%s' not found", e.Name)
}

type Resources []Resource

func (resources Resources) Lookup(name string) (Resource, bool) {
	for _, resource := range resources {
		if resource.Name() == name {
			return resource, true
		}
	}

	return nil, false
}

func (r *resource) ID() int              { return r.id }
func (r *resource) Name() string         { return r.name }
func (r *resource) PipelineID() int      { return r.pipelineID }
func (r *resource) PipelineName() string { return r.pipelineName }
func (r *resource) Type() string         { return r.type_ }
func (r *resource) Source() atc.Source   { return r.source }
func (r *resource) CheckEvery() string   { return r.checkEvery }
func (r *resource) Tags() atc.Tags       { return r.tags }
func (r *resource) CheckError() error    { return r.checkError }
func (r *resource) Paused() bool         { return r.paused }
func (r *resource) WebhookToken() string { return r.webhookToken }
func (r *resource) FailingToCheck() bool {
	return r.checkError != nil
}

func (r *resource) Reload() (bool, error) {
	row := resourcesQuery.Where(sq.Eq{"r.id": r.id}).
		RunWith(r.conn).
		QueryRow()

	err := scanResource(r, row)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *resource) Pause() error {
	_, err := psql.Update("resources").
		Set("paused", true).
		Where(sq.Eq{
			"id": r.id,
		}).
		RunWith(r.conn).
		Exec()

	return err
}

func (r *resource) Unpause() error {
	_, err := psql.Update("resources").
		Set("paused", false).
		Where(sq.Eq{
			"id": r.id,
		}).
		RunWith(r.conn).
		Exec()

	return err
}

func scanResource(r *resource, row scannable) error {
	var (
		configBlob      []byte
		checkErr, nonce sql.NullString
	)

	err := row.Scan(&r.id, &r.name, &configBlob, &checkErr, &r.paused, &r.pipelineID, &r.pipelineName, &nonce)
	if err != nil {
		return err
	}

	es := r.conn.EncryptionStrategy()

	var noncense *string
	if nonce.Valid {
		noncense = &nonce.String
	}

	decryptedConfig, err := es.Decrypt(string(configBlob), noncense)
	if err != nil {
		return err
	}

	var config atc.ResourceConfig
	err = json.Unmarshal(decryptedConfig, &config)
	if err != nil {
		return err
	}

	r.type_ = config.Type
	r.source = config.Source
	r.checkEvery = config.CheckEvery
	r.tags = config.Tags
	r.webhookToken = config.WebhookToken

	if checkErr.Valid {
		r.checkError = errors.New(checkErr.String)
	}

	return nil
}
