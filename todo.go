package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

type item struct {
	Task      string
	Done      bool
	CreatedAt time.Time
	CompletedAt time.Time
}

type Todos []item


func (t *Todos) Add( task string){

	todo := item{
		Task: task,
		Done: false,
		CreatedAt: time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)

}

func (t *Todos) Complete(index int) error{
	ls := *t

	if (index <= 0 || index > len(ls)){
		return errors.New("invalide index")
	}
	ls[index-1].CompletedAt = time.Now()
	ls[index-1].Done = true

	return nil
}


func (t *Todos) Delete(index int) error{
	ls := *t
	if (index <= 0 || index > len(ls)){
		return errors.New("invalide index")
	}

	*t = append(ls[:index-1],ls[index:]...)

	return nil
}


func (t *Todos) Load(filename string) error{

	file,err := ioutil.ReadFile(filename)
	if err != nil{
		if errors.Is(err, os.ErrNotExist){
			return nil	}
		return err
	}

	if len(file) == 0{
		return err
	}

	err = json.Unmarshal(file,t)

	if err != nil {
		return err
	}

	return nil
}


func (t *Todos) Store(filename string) error {

	data, err := json.Marshal(t)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data,0644)
}

func (t *Todos) Print(){


	table := simpletable.New()

	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Text: "No°"},
		{Align: simpletable.AlignCenter, Text: "Task"},
	{Align: simpletable.AlignCenter, Text: "Done"},
	{Align: simpletable.AlignCenter, Text: "Created At"},
	{Align: simpletable.AlignCenter, Text: "Complated At"},
		
}}
	var cells [][]*simpletable.Cell

	for i, item := range *t{
		i++
		task := blue(item.Task)
		done := red("No")
		if item.Done{
			task = green(fmt.Sprintf("\u2705 %s", task))
			done = green("Yes")
		}
		
		
		cells = append(cells,*&[]*simpletable.Cell{
			{Text: fmt.Sprintf("%d", i)},
			{Text :task},
			{Text: fmt.Sprintf("%s", done)},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}
	table.Body = &simpletable.Body{Cells:cells}
	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5,Text: red(fmt.Sprintf("%d Tasks are pending", t.CountTods()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()

}

func (t *Todos) CountTods() int{
	var s int = 0
	for _,item := range *t{
		if(!item.Done){
			s++
		}
	}
	return s
}