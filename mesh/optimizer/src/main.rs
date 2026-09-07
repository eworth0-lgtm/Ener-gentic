use tonic::{transport::Server, Request, Response, Status};

pub mod proto {
    tonic::include_proto!("mesh");
}

use proto::optimizer_server::{Optimizer, OptimizerServer};
use proto::{Forecast, Series, BatchForecast, BatchSeries};

extern "C" {
    fn simulate(u: f32, n: usize, tau: f32) -> f32;
}

#[derive(Debug, Default)]
pub struct OptimizerService {}

#[tonic::async_trait]
impl Optimizer for OptimizerService {
    async fn dispatch(
        &self,
        request: Request<Forecast>,
    ) -> Result<Response<Series>, Status> {
        let forecast = request.into_inner();
        let dispatch = optimize_dispatch(&forecast.values);
        Ok(Response::new(Series { values: dispatch }))
    }

    async fn dispatch_batch(
        &self,
        request: Request<BatchForecast>,
    ) -> Result<Response<BatchSeries>, Status> {
        let batch = request.into_inner();
        let mut series_list = Vec::new();
        
        for forecast in batch.forecasts {
            let dispatch = optimize_dispatch(&forecast.values);
            series_list.push(Series { values: dispatch });
        }
        
        Ok(Response::new(BatchSeries { series: series_list }))
    }
}

fn optimize_dispatch(forecast: &[f32]) -> Vec<f32> {
    let mut dispatch = Vec::with_capacity(forecast.len());
    
    for &demand in forecast {
        // Clamp and ramp
        let clamped = demand.max(0.0).min(100.0);
        
        // Call C++ sim for inverter lag
        let output = unsafe { simulate(clamped, 1, 0.05) };
        
        dispatch.push(output);
    }
    
    dispatch
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::]:50052".parse()?;
    let optimizer = OptimizerService::default();

    println!("Rust Optimizer running on {}", addr);

    Server::builder()
        .add_service(OptimizerServer::new(optimizer))
        .serve(addr)
        .await?;

    Ok(())
}
