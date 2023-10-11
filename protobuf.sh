protoc -I $PWD --go_out="./" $PWD/directive.proto
GEN_EAMS_LIBRARY=/Users/briananderson/Developer/Library/EmbeddedProto/protoc-gen-eams
protoc --plugin=protoc-gen-eams=$GEN_EAMS_LIBRARY -I $PWD --eams_out=/Users/briananderson/Arduino/kld7-sketch $PWD/directive.proto