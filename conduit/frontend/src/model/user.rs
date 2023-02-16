use std::ptr::null;

use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Clone, PartialEq,Serialize, Deserialize,Debug)]
#[serde(rename_all = "kebab-case")]
pub struct User{
    pub conduit_id: Uuid,
    pub external_auth_id:  uuid::Uuid,   
    pub conduit_display_name: String,
    pub user_type: String, 
    ///pub servers: Option<Vec<Uuid>>,
    pub status: Option<String>
}

// #[derive(Clone, PartialEq,Serialize, Deserialize,Debug)]
// pub struct Users {
//     pub users: Vec<User>
// }

impl Default for User {
    fn default () -> User {
        User{
            conduit_id: Uuid::default(),
            external_auth_id:Uuid::default(),
            conduit_display_name: String::new(),
            user_type: String::new(),
            //servers: None,
            status: None     
        
        }

    }
}

