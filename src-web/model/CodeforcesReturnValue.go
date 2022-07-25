package model

// InformationStruct
/*
# Return Information
ID := result.id
CID := result.contestId
PID := result.problem.index
PNAME := result.problem.name
CNAME := result.author.members.[name(Maybe NULL)/handle]
LANG := result.programmingLanguage

fileName := PID-PNAME-CNAME-LANG(CID#ID)
*/
type InformationStruct struct {
	ID    int64
	CID   int64
	PID   string
	PNAME string
	CNAME string
	LANG  string
	//TODO F: 对部分缩写进行重构
}
