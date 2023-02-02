

use gloo_net::http::Request;
use yew::prelude::*;
use yew_oauth2::prelude::*;



use gloo_net::{Error}; //add
use crate::model::user::User;

use crate::components::{card::Card}; 

#[function_component(JsonTest)]
pub fn app() -> Html {
    let users: UseStateHandle<Option<Vec<User>>> = use_state(|| None);
    let error: UseStateHandle<Option<Error>> = use_state(|| None);
    let auth = use_context::<OAuth2Context>();     
    let mut accessToken = String::new();
    if let Some(OAuth2Context::Authenticated (Authentication {access_token: authToken, ..})) = auth { 
        accessToken = authToken.clone();
    }else{
        accessToken = String::from("Not Logged In");
    }  
    {
        //create copies of states
        let users = users.clone();
        let error = error.clone();

        use_effect_with_deps(
            move |_| {
                wasm_bindgen_futures::spawn_local(async move {
                    let fetched_users = Request::get("http://localhost:9661/")
                    .header("Authentication", &accessToken) 
                    .send()
                    .await;
                    match fetched_users {
                        Ok(response) => {
                            let json = response.json::<Vec<User>>().await;
                            match json {
                                Ok(json_resp) => {
                                    users.set(Some(json_resp));
                                }
                                Err(e) => error.set(Some(e)),
                            }
                        }
                        Err(e) => error.set(Some(e)),
                    }
                });
                || ()
            },
            (),
        );
    }

    let user_list_logic = match users.as_ref() {
        Some(users) => users            
            .iter()
            .map(|user| {
                html! {
                  <Card user={user.clone() }/>
                }
            })
            .collect(),
        None => match error.as_ref() {
            Some(_) => {
                html! {
                     {"Error getting list of users"}
                }
            }
            None => {
                html! {
                  {"Loading"}
                }
            }
        },
    };

    html! {
      <>
        <h4> {"User"} </h4>
        {user_list_logic}
      </>
    }
}

