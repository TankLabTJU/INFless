%% 处理BATCH的原始数据
%qpsList.get(i)+","+throughtPerInstance*numberOfInstance+","+cpuCore*numberOfInstance+","+cpuCore*64*numberOfInstance+","+gpuCore*numberOfInstance+","+gpuCore*142*numberOfInstance+","+(cpuCore*64*numberOfInstance+gpuCore*142*numberOfInstance)+","+batchSize+","+throughtPerInstance+","+numberOfInstance+"\n"
%BATCH\resnet-50\3\resnet-50-350ms-300QPS-workload3-6s-60min.csv
 
data=temp;
resnet50_workload3_650msQPS_BATCH=zeros(sum(data(:,10)),4); %计算所有实例的数量
index=0;
for i=1:length(data);
    instanceNum=data(i,10);
    for j=1:instanceNum;
        index=index+1;
        resnet50_workload3_650msQPS_BATCH(index,:)=[data(i,8) data(i,3)/instanceNum data(i,5)/instanceNum data(i,9)];
    end
end