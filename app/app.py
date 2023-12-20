#!/usr/bin/python
from subprocess import check_output
from flask import Flask, render_template, stream_template, request
import logging

app = Flask(__name__)

script = ['awk', '-f', './renumber.awk']

@app.route('/', methods=['POST', 'GET'])
def renumber():
    if request.method == 'POST' and request.form['list'] != '':
        return stream_template('renumber.html',
                               renumbered=check_output(script,input=request.form['list'], encoding='utf8'))
    return render_template('renumber.html')
