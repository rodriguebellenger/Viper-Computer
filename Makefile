CC      = gcc
CFLAGS  = -Wall -Wextra -pedantic -O2
TARGET  = bytecode_interpreter/c
SRCS    = bytecode_interpreter/c.c
OBJS    = $(SRCS:.c=.o)

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) -o $(TARGET) $(OBJS)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

run: $(TARGET)
	./$(TARGET)

clean:
	rm -f $(TARGET) $(OBJS)
