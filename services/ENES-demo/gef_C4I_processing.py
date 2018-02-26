# -*- coding: latin-1 -*-
#  Copyright CERFACS (http://cerfacs.fr/)
#  Apache License, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0)
#
#  Author: Christian Page (2017)

import json
from icclim import *
import datetime
from icclim.util import callback
import time
import wget
import os
import ast
import io
import pdb

SAVEPATH = "/root/output/"
INPUT_PATH = "/root/input/"

try:
    to_unicode = unicode
except NameError:
    to_unicode = str

def netcdf_processing(params_file):

    # Extract parameters for processing
    with open(params_file) as data_file:    
        plist = json.load(data_file)

    # Extract processing function parameters
    var_name = plist["function"]["var_name"]
    out_var_name = plist["function"]["out_var_name"]
    slice_mode = plist["function"]["slice_mode"]
    if slice_mode == "None": slice_mode = None
    in_file = plist["function"]["in_file"]
    out_file = plist["function"]["out_file"]
    url = plist["function"]["url"]

    dd_b, mm_b, yyyy_b = map(int, plist["function"]["time_range_b"].split('-'))
    dd_e, mm_e, yyyy_e = map(int, plist["function"]["time_range_e"].split('-'))
        
    period = [datetime.datetime(yyyy_b,mm_b,dd_b), datetime.datetime(yyyy_e,mm_e,dd_e)]

    if plist["function"]["calc_operation"] == "time_avg":
        my_indice_params = {'indice_name': out_var_name, 
                            'calc_operation': 'mean',
                            }
    else:
        raise ValueError('Operation specified in calc_operation is not implemented: '+plist["function"]["calc_operation"])

    # Download file
    wget.download(url=url, out=in_file)

    # Check size
    size_file_in = os.path.getsize(in_file)

    # Add the dockerfile output path to the outfile name
    out_file = SAVEPATH+out_file

    # Launch processing using icclim
    try:
        if size_file_in>1.5:
            icclim.indice(user_indice=my_indice_params, in_files=in_file, var_name='tas', slice_mode=None, transfer_limit_Mbytes=2000, out_unit='days', out_file=out_file, callback=callback.defaultCallback2)
        else:
            icclim.indice(user_indice=my_indice_params, in_files=in_file, var_name='tas', slice_mode=None, out_unit='days', out_file=out_file, callback=callback.defaultCallback2)
    except:
        print(my_indice_params)
        print(in_file)
    os.remove(in_file)

def fromSting2Dict(textFile, jsonFile):

    try:
        with open(textFile, 'r') as myfile:
            Textfile = myfile.read().replace('\n', '')
        data = ast.literal_eval(Textfile)
    except:
        print("Textfile: ",Textfile)

    # Write JSON file
    with io.open(jsonFile, 'w', encoding='utf8') as outfile:
        str_ = json.dumps(data,
                          indent=4, sort_keys=True,
                          separators=(',', ': '), ensure_ascii=False)
        outfile.write(to_unicode(str_))

    return outfile

if __name__ == "__main__":

    input_string = INPUT_PATH+"processing_params.txt" 
    params_file = "processing_params.json"

    fromSting2Dict(input_string, params_file)
    netcdf_processing(params_file)
