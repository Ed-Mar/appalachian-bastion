
// //use log::log;
// //use serde_json::json;
// use yew::prelude::*;
// use yew_oauth2::prelude::*;

// use crate::model::user::User;

// //use crate::components::{card::Card}; //add


// use gloo_net::{http::Request, Error}; 

// // use uuid::Uuid;

// // use serde::Deserialize;





// #[function_component(BackendTest)]
// pub fn backend_test() -> Html {
    
//     //let error: UseStateHandle<Option<Error>> = use_state(|| None);
    
//     let mut accessToken=String::new();    
//     let auth = use_context::<OAuth2Context>();     
//     if let Some(OAuth2Context::Authenticated (Authentication {access_token: authToken, ..})) = auth { 
//         accessToken = authToken.clone();
//     }else{
//         accessToken = String::from("Not Logged In");
//     }   
        
//         ///let the_user: UseStateHandle<User> = use_state();        
        
//         let error = error.clone();    
// {   
//    // let the_user = the_user.clone();
//       //  let error = error.clone();
//         use_effect_with_deps(
//             move |_| {
//                 wasm_bindgen_futures::spawn_local(async move {
//                     let fetched_users = Request::get("http://localhost:9661/")
//                     .header("Authentication", &accessToken) 
//                     .send()
//                     .await;
//                     match fetched_users {
//                         Ok(response) => {
//                             let json = response.json::<User>().await;
//                             match json {
//                                 Ok(json_resp) => {    
//                                    // the_user.set(Some(json_resp));
                                                                   
//                                 }
//                                 Err(e) => error.set(Some(e)),

//                             }
//                         }
//                         Err(e) => error.set(Some(e)),
//                     }
//                 });
//                 || ()
//             },
//             (),
//         );
    
  
//             }

//     html! {
//       <>
//        /// <Card user={the_user.clone() }/>
        
        
//       </>
//     }
// }



// //     let auth = use_context::<OAuth2Context>();        
// //         if let Some(OAuth2Context::Authenticated (Authentication {access_token: authToken, ..})) = auth {                                  
// //             wasm_bindgen_futures::spawn_local( async move {                     
// //                 let s2 =  send_req(authToken.clone()).await;                            
// //             });
// //             html!(
// //                 <>{"somthing"}               

// //                 </>
// //             )
            
                
        
// //         } else {
// //             html!({"Fail"})
// //         }
        
    
// // }
// async fn send_req(access_token: String ) -> String {   
    
    
//     //let yeye = Request::new("http://localhost:9661/",);    
//     let user: User= Request::get("http://localhost:9661/")
//     .header("Authentication", &access_token) 
//     .send()
//     .await
//     .unwrap()
//     .json()
//     .await
//     .unwrap();
    
//     log::debug!("Response: '{:?}'", user);    
    
//     return String::from(String::from("Somthing"));
// }
// async fn new_user(access_token: String) -> String{
//     let req = Request::new("http://localhost:9661/",);
//     let resp= &req.method(gloo_net::http::Method::POST)
//     .header("Authentication", &access_token)    

//     .send()
//     .await
//     .unwrap();

//     return String::from(String::from("New User Method Run"));

// }