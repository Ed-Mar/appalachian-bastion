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
            <p class="conduit_id">{format! ("{}: {}","Conduit Id",user.conduit_id.clone())}</p>
            <p class="conduit_id">{format! ("{}: {}","External Authentication Id",user.external_auth_id.clone())}</p>
            <p class="conduit_id">{format! ("{}: {}","Conduit Display Name",user.conduit_display_name.clone())}</p>
            <p class="conduit_id">{format! ("{}: {}","User Type",user.user_type.clone())}</p>            
            <p class="status">{user.status.clone()}</p>
        </div>
    </div>
    }
}