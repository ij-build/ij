package util

import "sync"

func RunParallel(tasks ...func()) *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)

		go func(task func()) {
			defer wg.Done()
			task()
		}(task)
	}

	return wg
}

func RunParallelArgs(task func(string), args ...string) {
	tasks := []func(){}
	for _, arg := range args {
		tasks = append(tasks, capture(task, arg))
	}

	RunParallel(tasks...).Wait()
}

func capture(f func(string), arg string) func() {
	return func() {
		f(arg)
	}
}
