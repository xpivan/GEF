# -*- coding: latin-1 -*-
#  Copyright CERFACS (http://cerfacs.fr/)
#  Apache License, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0)
#
#  Author: Xavier Pivan (2017)

import json
from icclim import *
import datetime
from icclim.util import callback
import time
import ssl
import wget
import numpy as np
from glob import glob
from netCDF4 import Dataset
import pandas as pd
from matplotlib.mlab import griddata

DATA_INPUT = "/root/input/"
DATA_OUTPUT = "/root/output/"
SAVEFILE_NAME = "processing_dataset"
COLUMNS = ["date", "endpoint", "docker_image", "filename", "size_in", "size_out", "downloadTime", "processTime", "chunk_value"]
ENDPOINT = "carach5.ics.muni.cz"
DOCKER_IMAGE = "ubuntu"

def get_time():
    now_float = time.time()
    now_date = datetime.datetime.fromtimestamp(now_float).strftime('%Y-%m-%d %H:%M:%S')
    return now_float, now_date

def download_file(url, nc_filename):
    #Download the file in URL
    time_download_start, date_download_start = get_time()
    wget.download(url=url, out=nc_filename)
    time_download_end, date_download_end = get_time()
    time_download = time_download_end - time_download_start

    #Get data volume downloaded
    vol_data = os.path.getsize(nc_filename)

    return time_download, date_download_start, vol_data

def save_features(DATA_OUTPUT, SAVEFILE_NAME, data_proc):

    datafile = os.path.join(DATA_OUTPUT, SAVEFILE_NAME+".csv")
    newfile=True

    #Check if the file exists
    if os.path.exists(datafile):
        newfile = False

    df = pd.DataFrame(data_proc)
    df.columns = COLUMNS

    outf = open(datafile, "a")
    if newfile:
        df.to_csv(outf, index=False, columns=COLUMNS)
    else:
        df.to_csv(outf, index=False, columns=COLUMNS, header=False)
    outf.close()

def get_my_indice_params(nc_filename):

    my_indice_params = {'indice_name': 'tas', 
                            'calc_operation': 'mean',
                            }

    out_filename = nc_filename[0:-3]+'_'+my_indice_params['calc_operation']+nc_filename[-3::]

    return my_indice_params, out_filename

def netcdf_processing(my_indice_params, in_file, size_file_in, chunk_value, out_file):
    #Process starting time 
    time_process_start, date_process_start = get_time()

    # Launch processing using icclim
    size_file_in = size_file_in/1E9 # From Bytes to GBytes

    if size_file_in>1.5:
        icclim.indice(user_indice=my_indice_params, in_files=in_file, var_name='tas', slice_mode=None, transfer_limit_Mbytes=chunk_value, out_unit='days', out_file=out_file, callback=callback.defaultCallback2)
    else:
        icclim.indice(user_indice=my_indice_params, in_files=in_file, var_name='tas', slice_mode=None, out_unit='days', out_file=out_file, callback=callback.defaultCallback2)

    #Get data volume after processing
    vol_data_out = os.path.getsize(out_file)
    os.remove(in_file)

    #Process ending time
    time_process_end, date_process_end = get_time()
    process_time = time_process_end - time_process_start

    return process_time, date_process_start, vol_data_out

def create_netCDF(var, DATA_OUTPUT, output):

    #Instantiate a netcdf file
    dataset = Dataset(DATA_OUTPUT+output, 'w', format='NETCDF4_CLASSIC') 

    #Create the dimension
    lat = dataset.createDimension('lat', np.shape(var)[0])
    lon = dataset.createDimension('lon', np.shape(var)[1])

    #Create the variable
    lat = dataset.createVariable('lat', np.float32, ('lat',))
    lon = dataset.createVariable('lon', np.float32, ('lon',)) 

    #Create the temperature variable
    var = dataset.createVariable('tas', np.float32, ('lat', 'lon'))

    dataset.close()

def interp2d_mat(ncfile):

    #2 dimensional interpolation on the larger grid from the multiple dataset model
    fin = Dataset(ncfile)
    tas = fin.variables['tas'][0][:]
    lon = fin.variables['lon'][:]
    lat = fin.variables['lat'][:]
    fin.close()

    LON, LAT = np.meshgrid(lon,lat)
    LON1d = np.reshape(LON,np.size(LON))
    LAT1d = np.reshape(LAT,np.size(LAT))
    tas1d = np.reshape(tas[0,:,:],np.size(tas[0,:,:]))
    lon_new = np.linspace(lon[0],lon[-1],)
    lat_new = np.linspace(lat[0],lat[-1],)

    tas_new = griddata(LON1d,LAT1d,tas1d,lon_new,lat_new)

    return tas_new

if __name__ == "__main__":

    #We ignore the ssl warning which causes to abort the download
    ssl._create_default_https_context = ssl._create_unverified_context

    #A file text of selected URL we want to perform calculation on
    list_url = [line.rstrip('\n') for line in open(DATA_INPUT+'list_url.txt')]
    
    len_lu = len(list_url)
    save_file = SAVEFILE_NAME+"_"+str(len_lu)+"_files"

    if len_lu>1:
        #Array to store the resolution of each file
        check_res = np.zeros((len_lu+1,2))
    
    url_int = 0
    
    for url in list_url:
        #A hack to only keep the filename
        #nc_filename = url[71::]
        nc_filename = "tas_"+str(url_int+1)+".nc"
        chunk_value = 2000

        #Download file
        time_download, date_download_start, vol_data_in = download_file(url, nc_filename)

        #Extract instruction from json file
        my_indice_params, out_filename = get_my_indice_params(nc_filename)

        #Process the data using icclim library
        time_process, date_process_start, vol_data_out = netcdf_processing(my_indice_params, nc_filename, vol_data_in, chunk_value, out_filename)

        #data we want to save
        data_proc = [[date_download_start, ENDPOINT, DOCKER_IMAGE, out_filename, vol_data_in, vol_data_out, time_download, time_process, chunk_value]]
        save_features(DATA_OUTPUT, save_file, data_proc)

        #store spatial dimension
        if len_lu>1:
            fin = Dataset(out_filename)
            tas = fin.variables['tas'][:]
            lon = fin.variables['lon'][:]
            lat = fin.variables['lat'][:]
            check_res[url_int,:] = len(lat),len(lon)
            url_int +=1

    """if len_lu>1:

        new_ncfile = glob(DATA_OUTPUT+"*_mean.nc")

        nb_lat_max = np.int(np.max(check_res[:,0]))
        nb_lon_max = np.int(np.max(check_res[:,1]))

        tas_con = np.zeros((len_lu, nb_lat_max, nb_lon_max))
        ncfile_nb = 0

        for ncfile in new_ncfile:
            #2d interpolation on the larger grid from the multiple model dataset
            tas_con[ncfile_nb,:,:] = interp2d_mat(ncfile)
            ncfile_nb += 1

        tas_mean = np.mean(tas_con, axis=0)

        #We store the final result on a netcdf
        result = "average_of_"+str(len_lu)+"_files.nc"
        #create_netCDF(tas_mean, DATA_OUTPUT, result)"""
