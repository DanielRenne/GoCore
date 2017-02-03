package dbServices

type Mongo_Result_Repl_Conf struct {
	Config struct {
		Id       string                 `bson:"_id"`
		Version  int                    `bson:"version"`
		Members  []Mongo_Replica_Member `bson:"members"`
		Settings Mongo_Replica_Setting  `bson:"settings"`
	} `bson:"config"`
}

type Mongo_Replica_Member struct {
	Id          int    `bson:"_id"`
	Host        string `bson:"host"`
	ArbiterOnly bool   `bson:"arbiterOnly"`
	Hidden      bool   `bson:"hidden"`
	Priority    int    `bson:"priority"`
	Votes       int    `bson:"votes"`
	SlaveDelay  int    `bson:"slaveDelay"`
}

type Mongo_Replica_Setting struct {
	ChainingAllowed      bool `bson:"chainingAllowed"`
	HeartbeatTimeoutSecs int  `bson:"heartbeatTimeoutSecs"`
}
