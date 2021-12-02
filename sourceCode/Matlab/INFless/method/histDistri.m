function[] =histDistri(data,sum)
%生成一组随机数（正态分布） 
%绘制直方图
hist(data,sum); hold on;
%求出概率密度函数参数
[mu,sigma]=normfit(data);
%绘制概率密度函数
[n,x]=hist(data,sum);
y=normpdf(x,mu,sigma);
%处理一下数据，使得密度函数和最高点对齐
y=y/max(y)*max(n);
plot(x,y,'r-');
xlim([min(x),max(x)]);
end