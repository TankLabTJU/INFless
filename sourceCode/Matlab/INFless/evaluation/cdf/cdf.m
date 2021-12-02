function[p]=cdf(x,num)
xi = linspace(min(x),max(x),num);
F = ksdensity(x,xi,'function','cdf');
p=plot(xi,F);
xlim([min(x),max(x)]);
end