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
        user_email = request.form['email']
        user_password = request.form['password']

        # Validate form data
        if len(user_email) == 0 or len(user_password) == 0:
            # Form data failed validation; try again
            error = "Please supply both email and password"
        else:
            # Form data is valid; move along
            print(f'User set email: {user_email}.')
            print(f'User set password: {user_password}.')
            try:
                subprocess.check_call(
                    ['/home/pi/work/HomeGate/pi/setupEnvForDevice.sh', user_email, user_password])
            except Exception as e:
                error = 'failed to execute setupEnvForDevice.sh' + str(e)
                return render_template('setup.html', message=error)

            return redirect(url_for('thank_you'))

    # Render the index page
    return render_template('setup.html', message=error)
