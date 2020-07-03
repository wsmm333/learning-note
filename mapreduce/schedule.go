package mapreduce

import (
	"fmt"
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}

	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	 //构造任务
	 var wg sync.WaitGroup 
	 for i:= 0;i<ntasks;i++{
		wg.Add(1)    //正在运行的线程数加一

		//每个任务有自己的信息
		 var taskArgs DoTaskArgs
		 taskArgs.JobName = mr.jobName
		 taskArgs.Phase = phase
		 taskArgs.NumOtherPhase = nios
		 taskArgs.TaskNumber = i
		 if (phase == mapPhase) {
			 taskArgs.File = mr.files[i]
		 }
		 
		 go func(taskArgs DoTaskArgs){
			defer wg.Done()//当函数运行结束，减去一个进程数
			for{
				worker := <- mr.registerChannel
				if(call(worker,"Worker.DoTask",&taskArgs,nil) == false){   
					fmt.Printf("Worker %s PRC error!", worker)
					continue
				}else{
					go func(){
						mr.registerChannel <- worker      //回收worker
					}()
					break
				}
			}
			 
		 }(taskArgs)
	 }
	 wg.Wait() //当还有进程没有结束，堵塞状态
	debug("Schedule: %v phase done\n", phase)
}