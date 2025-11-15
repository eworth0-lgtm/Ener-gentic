#include <cmath>
#include <random>

extern "C" {
    float simulate(float u, size_t n, float tau) {
        // Inverter lag model: y = u * (1 - exp(-t/tau))
        // Add 0.1% jitter
        static std::random_device rd;
        static std::mt19937 gen(rd());
        static std::uniform_real_distribution<> dis(-0.001, 0.001);
        
        float t = 1.0f;  // time step
        float y = u * (1.0f - std::exp(-t / tau));
        
        // Add jitter
        y += y * dis(gen);
        
        return y;
    }
}
