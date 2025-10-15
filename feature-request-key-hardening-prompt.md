@internal/ @schema.sql @query.sql @sqlc.yaml @README.md can you help me implement this new feature for managing API keys:       
                                                                                                                                
the feature requirement is that each api_key should be associated with a project and allowed_origin (domain of the project that 
will use the api_key). It should also check if the api_key to be used is for dev only (is_dev flag bool), which if true should  
only allow allowed_origin: 'localhost'.                                                                                         
                                                                                                                                
- suggest a new schema for api_keys table, currently it is:                                                                     
                                                                                                                                
+--------------+--------------+------+-----+-------------------+-------------------+                                            
| Field        | Type         | Null | Key | Default           | Extra             |                                            
+--------------+--------------+------+-----+-------------------+-------------------+                                            
| api_key_id   | int          | NO   | PRI | NULL              | auto_increment    |                                            
| member_email | varchar(100) | NO   | MUL | NULL              |                   |                                            
| api_key_hash | varchar(255) | NO   |     | NULL              |                   |                                            
| created_at   | timestamp    | YES  |     | CURRENT_TIMESTAMP | DEFAULT_GENERATED |                                            
| expires_at   | timestamp    | YES  |     | NULL              |                   |                                            
+--------------+--------------+------+-----+-------------------+-------------------+                                            
                                                                                                                                
- i want the following logic and interface type for this api_key (see the details below):                                       
// api key object                                                                                                               
interface userAccessKey {                                                                                                       
  id: ...; // ikw bhla                                                                                                          
  key: string; // can be defined with crypto                                                                                    
  email: string;                                                                                                                
  project: string; // will help describe what it will be used for                                                               
  is_dev: boolean; // when is_dev is true, DO NOT allow if origin is not localhost, dev is responsible to not leak this         
  allowed_origin: string; // one url where they are allowed to call the api from. If is_dev = true, this should automatically be
localhost. If is_dev is false, and this is set, only allow from that url.                                                       
                                                                                                                                
// e.g.: if allowed_origin is set to links.app.dlsu-lscs.org, and the api key is used on a different url like                   
oms.app.dlsu-lscs.org, DO NOT allow and if possible have alerts for this                                                        
}                                                                                                                               
                                                                                                                                
- the /request-key should not be exposed to the public as i am planning to create a nextjs frontend dashboard for this and have 
google oauth in it, so thats the only way for the client to input their email for the /request-key endpoint (there should be no 
other way to hit the /request-key endpoint other than this frontend dashboard)                                                  
                                                                                                                                
- so right now, the logic that im thinking of for /request-key endpoints is:                                                    
  - the handler func for this should check for the new fields (is_dev and allowed_origin)                                       
    - if is_dev is true, then only allow 'localhost' for allowed_origin                                                         
    - if is_dev is false, then allow anything other than 'localhost' for allowed_origin                                         
    - there should only be one api_key per allowed_origin (so check for existing allowed_origin stored in database)             
      - for this is it better if one api_key per project? or per allowed_origin, as originally stated?                          
                                                                                                                                
help me design this properly, since this is mission critical that requests attention.
