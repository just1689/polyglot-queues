
docker service rm nsqlookupd
docker service rm nsqd
docker service rm nsqadmin

docker service create     --name nsqlookupd     --publish 4160:4160     --publish 4161:4161     --network my-services     nsqio/nsq /nsqlookupd
docker service create     --name nsqd     --publish 4150:4150     --publish 4151:4151     --network my-services     nsqio/nsq /nsqd --lookupd-tcp-address=nsqlookupd:4160
docker service create     --name nsqadmin     --publish 4171:4171     --network my-services     nsqio/nsq /nsqadmin --lookupd-http-address=nsqlookupd:4161
