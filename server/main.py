import zipfile
from flask import Flask, render_template, request, send_file, send_from_directory
from werkzeug import secure_filename
app = Flask(__name__)
import os

@app.route('/uploader', methods = ['GET', 'POST'])
def upload_file():
    if request.method == 'POST':
        if 'file' not in request.files:
            return 'no file sent', 400
        if 'user_name' not in request.form:
            return 'no user', 403
        f2 = open('files/' + request.form['user_name'] + '/pass', 'r')
        if request.form['pass'] == f2.read():
            f2.close()
            f = request.files['file']
            p = "files/"+request.form['user_name']+"/"+request.form['conf_name']+'/'+request.form['pretty_name']
            os.makedirs(p, exist_ok = True)
            f.save(p+'/file')
            f2 = open(p+'/config',"w+")
            f2.write(request.form['file_name'])
            f2.close()
            return 'wowo', 200
        f2.close()
        return 'no', 403

@app.route('/getfiles', methods = ['GET'])
def get_file():
    st = ""
    f2 = open('files/' + request.args['user_name'] + '/pass', 'r')
    if request.args['pass'] == f2.read():
        f2.close()
        for root, dirs, files in os.walk('files/'+request.args['user_name']+'/'+request.args['conf_name']):
            if len(dirs) > 0:
                for d in dirs:
                    st += request.args.get('user_name') + '/' + request.args.get('conf_name') + '/' + d + ' '
        print(st)
        return st, 200
    f2.close()
    return 'no', 403
    


@app.route('/file/<path:path>')
def send_js(path):
    return send_from_directory('files', path)

@app.route('/register', methods = ['GET'])
def register():
    path = 'files/' + request.args['user_name'] + '/pass'
    if not os.path.isfile(path):
        os.makedirs(path[:-4], exist_ok = True)
        f2 = open(path, 'w+')
        f2.write(request.args['pass'])
        f2.close()
    else:
        return "user already exists", 200

    return "done", 200

if __name__ == '__main__':
    app.run(debug = True)
