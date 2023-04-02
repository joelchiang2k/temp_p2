from flask import Flask, abort, request
import firebase_admin
from firebase_admin import credentials
from firebase_admin import db

cred = credentials.Certificate("serviceAccountKey.json")
firebase_admin.initialize_app(cred, {
    'databaseURL': 'https://registrytest-acbf4-default-rtdb.firebaseio.com/'
})

app = Flask(__name__)

@app.route('/package', methods=['POST'])
#POST package create
def post_package():
    auth = None
    auth_header = request.headers.get('X-Authorization')
    if not check_auth(auth_header):
        # add some auth function
        return abort(401, "Authentication failed (e.g. AuthenticationToken invalid or does not exist)")
    
    query = request.get_json()
    
    
    
    #read data from curl
    # open json and parse
    #if data not in db --> add the data
    #if data in db --> update the data
        #check what values it has
        #check if the url exists, check where the URL is
        
    
    
    #ingestion 
    
@app.route('/reset', methods=['DELETE'])
def reset_registry():
    auth_header = request.headers.get('X-Authorization')
    if not check_auth(auth_header):
        return abort(401, "You do not have permission to reset the registry.")
    
    ref = db.reference("/")
    ref.delete()
        
        
def check_auth(token):
    #add actual auth later
    if token is not None:
        return True

if __name__ == '__main__':
    app.run(debug=True)