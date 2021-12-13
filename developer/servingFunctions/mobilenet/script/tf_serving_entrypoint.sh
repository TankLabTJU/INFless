#!/bin/bash
sed -i "s/99/${BATCH_SIZE}/g" /home/app/function/batches/batching_parameters.txt
sed -i "s/5000000/${BATCH_TIMEOUT}/g" /home/app/function/batches/batching_parameters.txt
fwatchdog & tensorflow_model_server --rest_api_port=8501 --model_name=${MODEL_NAME} --model_base_path=${MODEL_BASE_PATH}/${MODEL_NAME} --enable_batching=true --batching_parameters_file=${BATCH_BASE_PATH}/batching_parameters.txt --per_process_gpu_memory_fraction=${GPU_MEM_FRACTION} "$@"
