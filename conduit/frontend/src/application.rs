use crate::components::*;
use yew::prelude::*;
use yew_nested_router::{components::*, prelude::*};
use yew_oauth2::prelude::*;

#[cfg(feature = "openid")]
use yew_oauth2::openid::*;

#[derive(Target, Debug, Clone, PartialEq, Eq)]
pub enum AppRoute {       
    Identity,
    #[target(index)]   
    Index,
}

#[function_component(Content)]
pub fn content() -> Html {
    let agent = use_auth_agent().expect("Requires OAuth2Context component in parent hierarchy");

    let login = {
        let agent = agent.clone();
        Callback::from(move |_: MouseEvent| {
            if let Err(err) = agent.start_login() {
                log::warn!("Failed to start login: {err}");
            }
        })
    };
    let logout = Callback::from(move |_: MouseEvent| {
        if let Err(err) = agent.logout() {
            log::warn!("Failed to logout: {err}");
        }
    });

    
    let openid_routes = html! (
        <li><Link<AppRoute> target={AppRoute::Identity}> { "Identity" } </Link<AppRoute>></li>
    );
    
    html!(
        <>
            <Router<AppRoute>>
                <Failure>
                    <ul>
                        <li><FailureMessage/></li>
                    </ul>
                </Failure>
                <Authenticated>
                    <p>
                        <button onclick={logout}>{ "Logout" }</button>
                    </p>                    
                    <ViewIdentity />
                    <Expiration/>                    
                </Authenticated>
                <NotAuthenticated>
                    <Switch<AppRoute> render={move |switch| match switch {
                        AppRoute::Index => html!(
                            <>
                                <p>
                                    { "You need to log in" }
                                </p>
                                <p>
                                    <button onclick={login.clone()}>{ "Login" }</button>
                                </p>
                            </>
                        ),
                        _ => html!(<LocationRedirect logout_href="/" />),
                    }} />
                </NotAuthenticated>
            </Router<AppRoute>>
        </>
    )
}

#[function_component(Application)]
pub fn app() -> Html {
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
            <h1> { "Login example "}</h1>

            <OAuth2
                {config}
                scopes={vec!["openid".to_string()]}
                >
                <Content/>
            </OAuth2>
        </>
    )
}



