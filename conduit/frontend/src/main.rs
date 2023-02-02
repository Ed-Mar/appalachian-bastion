#![recursion_limit = "1024"]

mod application;
mod components;
mod model;

use wasm_bindgen::prelude::*;

#[cfg(not(debug_assertions))]
const LOG_LEVEL: log::Level = log::Level::Info;
#[cfg(debug_assertions)]
const LOG_LEVEL: log::Level = log::Level::Trace;

pub fn main() -> Result<(), JsValue> {
    wasm_logger::init(wasm_logger::Config::new(LOG_LEVEL));
    log::info!("Starting application");
    yew::Renderer::<application::Application>::new().render();
    Ok(())
}
