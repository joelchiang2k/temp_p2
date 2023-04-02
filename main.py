from flask import Flask, abort, request, jsonify
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
    
    requestJSON = request.get_json()
    metadata = requestJSON['metadata']
    data = requestJSON['data']
    ID = metadata['ID']
    ref = db.reference("/")
    if ref.child(ID).get() is not None:
        return abort(409, "Package exists already.")
    else:
        ref.child(metadata['ID']).set({
            "name": metadata['Name'],
            "version": metadata['Version'],
            "ID": metadata['ID'],
            "content": data["Content"]
        })
    
    return jsonify({"metadata": metadata, "message": "Success. Check the ID in the returned metadata for the official ID"}),201 
    
@app.route('/package/byName/<name>', methods=['DELETE'])
def delete_package(name):
    print(name)
    ref = db.reference("/")
    table = ref.get()
    current_ID = None
    string = '"' + name + '"'
    for key, value in table.items():
      if value['name'] == name:
        current_ID = key
    if current_ID is None:
        return abort(404, "Package does not exist." + string)
    db.reference(current_ID).delete()
    return jsonify({"message": "Package is deleted."}), 200
        
    # if ref.child(name).get() is None:
    #     return abort(404, "Package does not exist.")
    # else:
    #     ref.child(name).delete()
    
    # return jsonify({"message": "Package is deleted."}), 200

@app.route('/reset', methods=['DELETE'])
def reset_registry():
    print("reset")
    auth_header = request.headers.get('X-Authorization')
    if not check_auth(auth_header):
        return abort(401, "You do not have permission to reset the registry.")
    
    ref = db.reference("/")
    ref.delete()
    
    return jsonify({"message": "Registry is reset."}), 200
        
def check_auth(token):
    #add actual auth later
    if token is not None:
        return True

if __name__ == '__main__':
    app.run(debug=True)