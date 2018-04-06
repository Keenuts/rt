
BIN=rt

all:
	cd src && go build -o ../$(BIN)

clean:
	$(RM) rt
