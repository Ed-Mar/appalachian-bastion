use yew::prelude::*;
use crate::model::user::User;

#[derive(Properties, PartialEq)]
pub struct CardProp {
    pub user: User
    
}

#[function_component(Card)]
pub fn card(CardProp { user }: &CardProp) -> Html {
    html! {
    <div class="m-3 p-4 border rounded d-flex align-items-center">
        
        <div class="">            
            <p class="fw-normal mb-1">{user.conduit_id.clone()}</p>
            <p class="fw-normal mb-1">{user.external_auth_id.clone()}</p>
            <p class="fw-normal mb-1">{user.external_auth_provider.clone()}</p>
            <p class="fw-normal mb-1">{user.external_auth_client_id.clone()}</p>
            <p class="fw-normal mb-1">{user.external_user_name.clone()}</p>
            <p class="fw-normal mb-1">{user.display_user_name.clone()}</p>
            <p class="fw-normal mb-1">{user.user_type.clone()}</p>
           //<p class="fw-normal mb-1">{user.servers.clone()}</p>
            <p class="fw-normal mb-1">{user.status.clone()}</p>
        </div>
    </div>
    }
}