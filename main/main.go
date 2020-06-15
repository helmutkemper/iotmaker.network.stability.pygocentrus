package main

import (
	//"github.com/helmutkemper/iotmaker.docker.util.whaleAquarium/factoryContainerMongoDB"
	//"github.com/helmutkemper/iotmaker.docker.util.whaleAquarium/factoryWhaleAquarium"
	pygocentrus "github.com/helmutkemper/iotmaker.network.stability.pygocentrus"
	"time"
)

func main() {
	//var id string
	var err error

	/*ch := factoryWhaleAquarium.NewPullStatusMonitor()

	  err, id = factoryContainerMongoDB.NewSingleEphemeralInstanceMongo(
	    "mongoEphemeral",
	    "test",
	    factoryContainerMongoDB.KMongoDBVersionTag_latest,
	    ch,
	  )
	  if err != nil {
	    panic(err)
	  }
	  _ = id
	*/

	l := pygocentrus.Listen{
		In: pygocentrus.Connection{
			Address:  "127.0.0.1:27017",
			Protocol: "tcp",
		},
		Out: pygocentrus.Connection{
			Address:  "127.0.0.1:27016",
			Protocol: "tcp",
		},
		Pygocentrus: pygocentrus.Pygocentrus{
			Enabled: true,
			Delay: pygocentrus.RateMaxMin{
				Rate: 0.5,
				Min:  int(time.Millisecond * 300),
				Max:  int(time.Millisecond * 600),
			},
			DontRespond: pygocentrus.RateMaxMin{
				Rate: 0,
				Min:  0,
				Max:  0,
			},
			ChangeLength: 0,
			ChangeContent: pygocentrus.ChangeContent{
				ChangeRateMin:  0,
				ChangeRateMax:  0,
				ChangeBytesMin: 0,
				ChangeBytesMax: 0,
				Rate:           0,
			},
			DeleteContent: 0,
		},
	}
	err = l.Listen()
	if err != nil {
		panic(err)
	}
}
