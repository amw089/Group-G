using System;
using System.Security.Cryptography;
using System.Collections.Generic;
using System.Text;
using System.Linq;

namespace myapp
{
    class Program
    {  
        public static readonly Guid OURnamespace = new Guid("d9b2d63d-a233-4123-847a-76838bf2413a");
        static StringBuilder sb = new StringBuilder();        
        public static Guid uuidV5(string name) {
                byte[] nameBytes = Encoding.UTF8.GetBytes(name);
                byte[] namespaceBytes = OURnamespace.ToByteArray();
                SwapByteOrder(namespaceBytes);
                byte[] hash;

                var incrementalHash = IncrementalHash.CreateHash(HashAlgorithmName.SHA1);
                incrementalHash.AppendData(namespaceBytes);
			    incrementalHash.AppendData(nameBytes);
			    hash = incrementalHash.GetHashAndReset();
                
                // most bytes from the hash are copied straight to the bytes of the new GUID (steps 5-7, 9, 11-12)
                byte[] newGuid = new byte[16];
                Array.Copy(hash, 0, newGuid, 0, 16);
                // set the four most significant bits (bits 12 through 15) of the time_hi_and_version field to the appropriate 4-bit version number from Section 4.1.3 (step 8)
                newGuid[6] = (byte)((newGuid[6] & 0x0F) | (5 << 4));
                // set the two most significant bits (bits 6 and 7) of the clock_seq_hi_and_reserved to zero and one, respectively (step 10)
                newGuid[8] = (byte)((newGuid[8] & 0x3F) | 0x80);
                // convert the resulting UUID to local byte order (step 13)
                SwapByteOrder(newGuid);
		        Guid guidV5 = new Guid(newGuid);

                return guidV5;

        }
        public static Dictionary<string, string>  populateNames() {

            // Create Dictionary.
            Dictionary<string, string> hash = new Dictionary<string, string>();
            
            string line = "";
            string key = "";  

            // Read the file and display it line by line.  
            System.IO.StreamReader file = new System.IO.StreamReader("names.txt");  

            while((line = file.ReadLine()) != null)  
            {  
                key = uuidV5(line).ToString();
                try {
                    hash.Add(key, line);
                }
                catch(System.ArgumentException e) {
                    sb.Append(e.Message);
                }
                
            }  
            
            file.Close(); 

            return hash;
        }

        public static Dictionary<string, string>  populateDictionary() {

            // Create Dictionary.
            Dictionary<string, string> dict = new Dictionary<string, string>();
            
            string line = "";
            string key = "";  

            // Read the file and display it line by line.  
            System.IO.StreamReader file = new System.IO.StreamReader("dictionary.txt");  

            while((line = file.ReadLine()) != null)  
            {
                SHA256 sha256Hash = SHA256.Create();
                byte[] bytes = sha256Hash.ComputeHash(Encoding.UTF8.GetBytes(line));
                // Convert byte array to a string   
                key = BitConverter.ToString(bytes).Replace("-", string.Empty);

                try {
                    dict.Add(key, line);
                }
                catch(System.ArgumentException e) {
                    sb.Append(e.Message);
                }
                
            }  
            
            file.Close(); 

            return dict;
        }

        static void Main(string[] args)
        {
            Dictionary<string, string>  nameMap = populateNames();
            Dictionary<string, string>  dictMap = populateDictionary();
             
            string line = "";

            // Read the file and display it line by line.  
            System.IO.StreamReader file = new System.IO.StreamReader("database_dump.csv");  
            System.IO.StreamWriter logFile = new System.IO.StreamWriter("log");

            while((line = file.ReadLine()) != null)  
            {  
                
                String[] separator = { "," }; 
                String[] strlist = line.Split(separator,StringSplitOptions.RemoveEmptyEntries); 
                if(strlist[0] == "username") {
                    System.Console.WriteLine("username,password,last_access");
                }
                else {
                    System.Console.WriteLine(nameMap[strlist[0]]+","+dictMap[strlist[1]]+","+ decodeTimestamp(strlist[2]));
                }
                
                
            }  
            
            file.Close();
            logFile.Write(sb.ToString()+"\n");
            logFile.Close(); 
        }

        private static string decodeTimestamp(string unixT) {
                double timestamp = Convert.ToDouble(unixT);
                System.DateTime dateTime = new System.DateTime(1970, 1, 1, 0, 0, 0, 0);
                dateTime = dateTime.AddSeconds(timestamp);
                
                TimeZoneInfo cstZone = TimeZoneInfo.FindSystemTimeZoneById("America/Regina");
                DateTime cstTime = TimeZoneInfo.ConvertTimeFromUtc(dateTime, cstZone);

                string formats = "yyyy-MM-ddTHH:mm:ss";
                return cstTime.ToString(formats)+"-0600";
        } 

        private static void SwapByteOrder(byte[] guid) {
		SwapBytes(guid, 0, 3);
		SwapBytes(guid, 1, 2);
		SwapBytes(guid, 4, 5);
		SwapBytes(guid, 6, 7);
	}

        private static void SwapBytes(byte[] guid, int left, int right) {
            byte temp = guid[left];
            guid[left] = guid[right];
            guid[right] = temp;
        }
    }
}
