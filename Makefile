gfuzz-inc:
	./GFuzz/scripts/fuzz-mount.sh /mnt/c/Users/jerem/Documents/GitHub/go-fuzzing-seminar/examples/increment/ /mnt/c/Users/jerem/Documents/GitHub/go-fuzzing-seminar/out

gfuzz-hello:
	./GFuzz/scripts/fuzz-mount.sh /mnt/c/Users/jerem/Documents/GitHub/go-fuzzing-seminar/examples/hello/ /mnt/c/Users/jerem/Documents/GitHub/go-fuzzing-seminar/out

gfuzz-conc:
	./GFuzz/scripts/fuzz-mount.sh /mnt/c/Users/jerem/Documents/GitHub/go-fuzzing-seminar/examples/conc/ /mnt/c/Users/jerem/Documents/GitHub/go-fuzzing-seminar/out

clean-out:
	rm -rf ./out/*