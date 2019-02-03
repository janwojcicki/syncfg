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
        f = request.files['file']
        p = "files/"+request.form['user_name']+"/"+request.form['conf_name']+'/'+request.form['pretty_name']
        os.makedirs(p, exist_ok = True)
        f.save(p+'/file')
        f2 = open(p+'/config',"w+")
        f2.write(request.form['file_name'])
        f2.close()
        return 'wowo', 200

@app.route('/getfiles', methods = ['GET'])
def get_file():
    st = ""
    print(request.args)
    for root, dirs, files in os.walk('files/'+request.args['user_name']+'/'+request.args['conf_name']):
        if len(dirs) > 0:
            for d in dirs:
                st += request.args.get('user_name') + '/' + request.args.get('conf_name') + '/' + d + ' '
    print(st)
    return st


@app.route('/file/<path:path>')
def send_js(path):
    return send_from_directory('files', path)

if __name__ == '__main__':
    app.run(debug = True)
