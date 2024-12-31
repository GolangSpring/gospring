package main

import (
	"go-spring/application"
	"go-spring/domains/ansible"
	"go-spring/pkg/postgres"
)

func main() {
	//play :=
	//	repository.Play{
	//		Name: "test",
	//		Variables: map[string]any{
	//			"example_var1": "test",
	//		},
	//		Hosts: "test",
	//		Tasks: []*repository.Task{
	//			{
	//				Name:         "task1",
	//				ModuleAction: &action.CommandModuleAction{Command: "echo {{ example_var1}}"},
	//			},
	//		},
	//	}
	//book := repository.Playbook{
	//	ID:   uuid.New().String(),
	//	Name: "Test",
	//	Plays: []repository.Play{
	//		play,
	//	}}
	//
	//_inventory := repository.Inventory{ID: uuid.New().String(), HostGroups: map[string][]string{
	//	"test": {"localhost", "127.0.0.1"},
	//}}

	//bookExecutor := service.NewPlaybookExecutor()
	//
	//go func() {
	//	for {
	//		select {
	//		case out := <-bookExecutor.StdOutChannel:
	//			fmt.Println("STDOUT:", out)
	//		case err := <-bookExecutor.StdErrChannel:
	//			fmt.Println("STDERR:", err)
	//		}
	//	}
	//}()
	//
	//if err := bookExecutor.ExecutePlaybook(ansible.TODO(), &book, &_inventory); err != nil {
	//	fmt.Println(err)
	//}
	configPath := "/Users/william_w_chen/GolandProjects/go-spring/cmd/api/config.yaml"
	appConfig := application.MustNewAppConfig(configPath)
	postgresConfig := postgres.MustNewPostgresDataSourceConfig(configPath)
	postgresContext := postgres.MustNewPostgresApplicationContext(postgresConfig)
	app := application.MustNewApplication(appConfig)
	ansibleAppContext := ansible.MustNewAnsibleAppContext()
	app.InjectContextCollection(postgresContext, ansibleAppContext)
	app.Run()
}
