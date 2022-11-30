use crate::components::*;
use yew::prelude::*;
use yew_oauth2::prelude::*;

use yew_oauth2::openid::*;
use yew_oauth2::openid::Client;


#[derive(Clone, Default, Debug, PartialEq, Properties)]
pub struct Props {}

pub struct Application {}

impl Component for Application {
    type Message = ();
    type Properties = Props;

    fn create(_ctx: &Context<Self>) -> Self {
        Self {}
    }

    fn view(&self, ctx: &Context<Self>) -> Html {
        let login = ctx.link().callback_once(|_: MouseEvent| {
            OAuth2Dispatcher::<Client>::new().start_login();
        });
        let logout = ctx.link().callback_once(|_: MouseEvent| {
            OAuth2Dispatcher::<Client>::new().logout();
        });

        #[cfg(feature = "openid")]
        let config = Config {
            client_id: "dev-conduit-rust".into(),            
            issuer_url: "http://keycloak.test/realms/gatehouse".into(),
            additional: Additional {
                
                after_logout_url:Some("http://localhost:8080".into()),
                ..Default::default()
            }
        };

        html!(
            <>
             <OAuth2
                {config}
                scopes={vec!["openid".to_string()]}
                >
                <Failure>
                    <div>
                        <FailureMessage/>
                    </div>
                </Failure>
                <Authenticated>
                    <div class="header">
                        <div class="header-item"> <button onclick={logout}>{ "Logout" }</button> </div>
                        <div class="header-item"> { " | Status: Logged In | "} </div>                        
                    </div>
                    <BackendTest />
                    <ViewIdentity />
                </Authenticated>
                <NotAuthenticated>

                            <button onclick={login.clone()}>{ "Login" }</button>

                </NotAuthenticated>



            </OAuth2>
            </>
        )
    }
}
