clc;
clear;
load workload3_190QPS;
data=floor(workload3_190QPS*50/190);

SLO=200;
execute=69.8;
timeout=SLO-execute;
request=0;
for i=1:length(data);
    interval=1000/data(i);
    if timeout<interval;
        request=request+data(i);
        data(i,2)=data(i,1);
    else 
        concurrency=floor(timeout/interval);
        if concurrency>4;
            concurrency=4;
        end
        request=request+ceil(data(i,1)/concurrency);
        data(i,2)=ceil(data(i,1)/concurrency);
    end
end
        

execute=168;  
timeout=SLO-execute;
request=0;
for i=1:length(data);
    interval=1000/data(i);
    if timeout<interval;
        request=request+1;
        data(i,3)=data(i);
    else 
        concurrency=floor(timeout/interval);
        if concurrency>16;
            concurrency=16;
        end
        request=request+ceil(data(i)/concurrency);
        data(i,3)=ceil(data(i)/concurrency);
    end
end
        

%% resnet-50
% clc;
% clear;
% load workload3_190QPS; 
% data=workload3_190QPS;
% 
% SLO=500; 
% execute=433;
% timeout=SLO-execute;
% request=0;
% for i=1:length(data);
%     interval=1000/data(i);
%     if timeout<interval;
%         request=request+data(i);
%         data(i,2)=data(i,1);
%     else 
%         concurrency=floor(timeout/interval);
%         if concurrency>4;
%             concurrency=4;
%         end
%         request=request+ceil(data(i,1)/concurrency);
%         data(i,2)=ceil(data(i,1)/concurrency);
%     end
% end
%         
% 
% execute=112;
% timeout=SLO-execute;
% request=0;
% for i=1:length(data);
%     interval=1000/data(i);
%     if timeout<interval;
%         request=request+1;
%         data(i,3)=data(i);
%     else 
%         concurrency=floor(timeout/interval);
%         if concurrency>4;
%             concurrency=4;
%         end
%         request=request+ceil(data(i)/concurrency);
%         data(i,3)=ceil(data(i)/concurrency);
%     end
% end
%         