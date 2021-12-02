package scs.util.format;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.IOException;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import javax.imageio.stream.FileImageInputStream;

import sun.misc.BASE64Decoder;
import sun.misc.BASE64Encoder;

@SuppressWarnings("restriction")
public class Base64Util{
    /**
     * 字符串转图片
     * @param base64Str
     * @return
     */
    public static byte[] decode(String base64Str){
        byte[] b = null;
        BASE64Decoder decoder = new BASE64Decoder();
        try {
            b = decoder.decodeBuffer(replaceEnter(base64Str));
        } catch (IOException e) {
            e.printStackTrace();
        }
        return b;
    }
    
    /**
     * 图片转字符串
     * @param image
     * @return
     */
    public static String encode(byte[] image){
        BASE64Encoder decoder = new BASE64Encoder();
        return replaceEnter(decoder.encode(image));
    }
    
    public static String encode(String uri){
        BASE64Encoder encoder = new BASE64Encoder();
        return replaceEnter(encoder.encode(uri.getBytes()));
    }
    
    /**
     * 
     * @path    图片路径
     * @return
     */
    
    public static byte[] imageTobyte(String path){
        byte[] data = null;
        FileImageInputStream input = null;
        try {
            input = new FileImageInputStream(new File(path));
            ByteArrayOutputStream output = new ByteArrayOutputStream();
            byte[] buf = new byte[1024];
            int numBytesRead = 0;
            while((numBytesRead = input.read(buf)) != -1){
                output.write(buf, 0, numBytesRead);
            }
            data = output.toByteArray();
            output.close();
            input.close();
            
        } catch (Exception e) {
            e.printStackTrace();
        }
        
        return data;
    }
    
    
    
    public static String replaceEnter(String str){
        String reg ="[\n-\r]";
        Pattern p = Pattern.compile(reg);
        Matcher m = p.matcher(str);
        return m.replaceAll("");
    }
    public static void main(String[] args){
    	System.out.println(Base64Util.encode(Base64Util.imageTobyte(Base64Util.class.getResource("/").getPath()+"conf/cat.jpg")));
    }
    
}