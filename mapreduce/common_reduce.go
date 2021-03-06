package mapreduce

import (
	"fmt"
	"encoding/json"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	// You can find the intermediate file for this reduce task from map task number
	// m using reduceName(jobName, m, reduceTaskNumber).
	// Remember that you've encoded the values in the intermediate files, so you
	// will need to decode them. If you chose to use JSON, you can read out
	// multiple decoded values by creating a decoder, and then repeatedly calling
	// .Decode() on it until Decode() returns an error.
	//
	// You should write the reduced output in as JSON encoded KeyValue
	// objects to a file named mergeName(jobName, reduceTaskNumber). We require
	// you to use JSON here because that is what the merger than combines the
	// output from all the reduce tasks expects. There is nothing "special" about
	// JSON -- it is just the marshalling format we chose to use. It will look
	// something like this:
	//
	// enc := json.NewEncoder(mergeFile)
	// for key in ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()

	//outputFile,err := os.OpenFile(outFile,os.O_WRONLY|os.O_CREATE,0666) //读入中间文件??
    //defer outputFile.Close()
    /*if err != nil {
        fmt.Println("err")
        return
    }*/
    res :=make(map[string] []string)
    for i := 0;i<nMap;i++{
        inputFile,err := os.Open(reduceName(jobName,i,reduceTaskNumber))
        defer inputFile.Close()
        if err != nil{
            fmt.Println("error")
            return
        }
        dec := json.NewDecoder(inputFile)   
        for{
            var kv KeyValue
            err := dec.Decode(&kv)     //读取文件中kv值
            if err!=nil {
                break
            }
            res[kv.Key] = append(res[kv.Key],kv.Value)  //向切片中追加数据    //相同key的kv合并，res结构：map[string] []string
        }
    }
    keys := make([]string,0)
    for k,_ :=range res{
        keys = append(keys,k)
    }
	sort.Strings(keys)
	
	mergeFileName := mergeName(jobName, reduceTaskNumber)    //在common里
	file, err := os.Create(mergeFileName)       //创建文件，如果文件已存在，会将文件清空。
	defer file.Close()
	if err != nil {
		fmt.Println("Cannot open file ", mergeFileName)
	}
 
    enc := json.NewEncoder(file)
    for _,key := range keys{
        enc.Encode(KeyValue{key,reduceF(key,res[key])}) //对key的次数进行计数
    }
    //file , _ := ioutil.ReadFile(outFile)
    //fmt.Println(string(file))
}
