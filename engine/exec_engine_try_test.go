package engine_test

import (
	"github.com/concourse/atc"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/engine"
	"github.com/concourse/atc/engine/fakes"
	"github.com/concourse/atc/event"
	"github.com/concourse/atc/exec"
	"github.com/concourse/atc/worker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"

	execfakes "github.com/concourse/atc/exec/fakes"
)

var _ = Describe("Exec Engine with Try", func() {

	var (
		fakeFactory         *execfakes.FakeFactory
		fakeDelegateFactory *fakes.FakeBuildDelegateFactory
		fakeDB              *fakes.FakeEngineDB

		execEngine engine.Engine

		buildModel db.Build
		logger     *lagertest.TestLogger

		fakeDelegate *fakes.FakeBuildDelegate
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("test")

		fakeFactory = new(execfakes.FakeFactory)
		fakeDelegateFactory = new(fakes.FakeBuildDelegateFactory)
		fakeDB = new(fakes.FakeEngineDB)

		execEngine = engine.NewExecEngine(fakeFactory, fakeDelegateFactory, fakeDB)

		fakeDelegate = new(fakes.FakeBuildDelegate)
		fakeDelegateFactory.DelegateReturns(fakeDelegate)

		buildModel = db.Build{ID: 84}
	})

	Context("running try steps", func() {
		var (
			taskStepFactory *execfakes.FakeStepFactory
			taskStep        *execfakes.FakeStep

			inputStepFactory *execfakes.FakeStepFactory
			inputStep        *execfakes.FakeStep
		)

		BeforeEach(func() {
			taskStepFactory = new(execfakes.FakeStepFactory)
			taskStep = new(execfakes.FakeStep)
			taskStep.ResultStub = successResult(true)
			taskStepFactory.UsingReturns(taskStep)
			fakeFactory.TaskReturns(taskStepFactory)

			inputStepFactory = new(execfakes.FakeStepFactory)
			inputStep = new(execfakes.FakeStep)
			inputStep.ResultStub = successResult(true)
			inputStepFactory.UsingReturns(inputStep)
			fakeFactory.GetReturns(inputStepFactory)
		})

		Context("constructing steps", func() {
			var (
				fakeDelegate          *fakes.FakeBuildDelegate
				fakeInputDelegate     *execfakes.FakeGetDelegate
				fakeExecutionDelegate *execfakes.FakeTaskDelegate
			)

			BeforeEach(func() {
				fakeDelegate = new(fakes.FakeBuildDelegate)
				fakeDelegateFactory.DelegateReturns(fakeDelegate)

				fakeInputDelegate = new(execfakes.FakeGetDelegate)
				fakeDelegate.InputDelegateReturns(fakeInputDelegate)

				fakeExecutionDelegate = new(execfakes.FakeTaskDelegate)
				fakeDelegate.ExecutionDelegateReturns(fakeExecutionDelegate)

				plan := atc.Plan{
					Try: &atc.TryPlan{
						Step: atc.Plan{
							Get: &atc.GetPlan{
								Name: "some-input",
							},
						},
					},
					Task: &atc.TaskPlan{
						Name: "some task",
					},
				}

				build, err := execEngine.CreateBuild(buildModel, plan)
				Ω(err).ShouldNot(HaveOccurred())
				build.Resume(logger)
			})

			It("constructs the step correctly", func() {
				Ω(fakeFactory.GetCallCount()).Should(Equal(1))
				sourceName, workerID, delegate, _, _, _, _ := fakeFactory.GetArgsForCall(0)
				Ω(sourceName).Should(Equal(exec.SourceName("some-input")))
				Ω(workerID).Should(Equal(worker.Identifier{
					BuildID:      84,
					Type:         worker.ContainerTypeGet,
					Name:         "some-input",
					StepLocation: 2,
				}))

				Ω(delegate).Should(Equal(fakeInputDelegate))
				_, _, location, hook := fakeDelegate.InputDelegateArgsForCall(0)
				Ω(location).Should(Equal(event.OriginLocation{
					ParentID:      0,
					ID:            2,
					ParallelGroup: 0,
				}))
				Ω(hook).Should(Equal(""))
			})
		})

		Context("when the inner step fails", func() {
			BeforeEach(func() {
				inputStep.ResultStub = successResult(false)
			})

			It("runs the next step", func() {
				plan := atc.Plan{
					HookedCompose: &atc.HookedComposePlan{
						Step: atc.Plan{
							Try: &atc.TryPlan{
								Step: atc.Plan{
									Get: &atc.GetPlan{
										Name: "some-input",
									},
								},
							},
						},
						Next: atc.Plan{
							Task: &atc.TaskPlan{
								Name:   "some-resource",
								Config: &atc.TaskConfig{},
							},
						},
					},
				}

				build, err := execEngine.CreateBuild(buildModel, plan)

				Ω(err).ShouldNot(HaveOccurred())

				build.Resume(logger)

				Ω(inputStep.RunCallCount()).Should(Equal(1))
				Ω(inputStep.ReleaseCallCount()).Should(Equal(3))

				Ω(taskStep.RunCallCount()).Should(Equal(1))
				Ω(taskStep.ReleaseCallCount()).Should(Equal(1))
			})
		})
	})
})
