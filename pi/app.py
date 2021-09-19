#!/usr/bin/python3

import subprocess
from flask import Flask, render_template, request, redirect, url_for


app = Flask(__name__)


@app.route('/')
def index():
    return render_template('index.html')


@app.route('/thank-you')
def thank_you():
    return render_template('thank-you.html')


@app.route('/setup', methods=['GET', 'POST'])
def setup():
    error = ""
    if request.method == 'POST':
        # Form being submitted; grab data from form.
        gate_id = request.form['gate_id']

        # Validate form data
        if len(gate_id) == 0:
            # Form data failed validation; try again
            error = "Please supply a valid gate ID"
        else:
            # Form data is valid; move along
            print('User set gate ID:' + gate_id)
            try:
                subprocess.check_call(
                    ['/home/pi/HomeGate/pi/setupEnvForDevice.sh', gate_id])
            except Exception as e:
                error = 'failed to execute setupEnvForDevice.sh' + str(e)
                return render_template('setup.html', message=error)

            return redirect(url_for('thank_you'))

    # Render the index page
    return render_template('setup.html', message=error)
