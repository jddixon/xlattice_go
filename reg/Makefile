p.pb.go: p.proto
	mkdir -p _pb
	protoc --go_out=_pb $<
	cat _pb/$@\
	|gofmt >$@
	rm -rf _pb

#r.pb.go: r.proto
#	mkdir -p _pb
#	protoc --go_out=_pb $<
#	cat _pb/$@\
#	|gofmt >$@
#	rm -rf _pb
