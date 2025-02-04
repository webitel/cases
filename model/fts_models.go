package model

type FtsCase struct {
	Id            string `json:"id,omitempty"`
	Description   string `json:"description,omitempty"`
	CloseResult   string `json:"closeResult,omitempty"`
	RatingComment string `json:"ratingComment,omitempty"`
}

type FtsCaseComment struct {
	ParentId int64  `json:"parentId,omitempty"`
	Comment  string `json:"comment,omitempty"`
}
