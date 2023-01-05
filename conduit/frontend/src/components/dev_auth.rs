
use yew::prelude::*;
use yew_oauth2::prelude::*;

use gloo_net::http::Request;



#[function_component(BackendTest)]
pub fn backend_test() -> Html {
    let auth = use_context::<OAuth2Context>();        
        if let Some(OAuth2Context::Authenticated (Authentication {access_token: authToken, ..})) = auth {            
                      
            wasm_bindgen_futures::spawn_local( async move {                     
                let s2 =  send_req(authToken.clone()).await;           
                 
            });
            

            html!(
                <>
                    {"somthing"}               

                </>
            )
            
                
        
        } else {
            html!({"Fail"})
        }
        
    
}
async fn send_req(access_token: String ) -> String {   
    
    let yeye = Request::new("http://localhost:9666/",);    
    let resp =&yeye.method(gloo_net::http::Method::POST)
    .header("Authentication", &access_token)    
    
    .send()
    .await
    .unwrap();

    return String::from(String::from("Somthing"));
}
async fn new_user(access_token: String) -> String{
    let req = Request::new("http://localhost:9661/",);
    let resp= &req.method(gloo_net::http::Method::POST)
    .header("Authentication", &access_token)    

    .send()
    .await
    .unwrap();

    return String::from(String::from("New User Method Run"));

}