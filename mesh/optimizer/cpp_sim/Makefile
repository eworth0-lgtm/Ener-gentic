CXX = g++
CXXFLAGS = -O3 -fPIC -std=c++17
TARGET = libsim.a

all: $(TARGET)

$(TARGET): sim.o
	ar rcs $(TARGET) sim.o

sim.o: sim.cpp
	$(CXX) $(CXXFLAGS) -c sim.cpp -o sim.o

clean:
	rm -f *.o $(TARGET)

.PHONY: all clean
