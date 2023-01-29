using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;
using System.Text.Json;

using TinyCsvParser;

namespace SelectionParseUtility
{
    class Program
    {
        const string CONFIG_FILE = "appSettings.json";
        static string _inputFile;

        static void Main(string[] args)
        {
            if (!ParseConfiguration())
            {
                Console.WriteLine("Error parsing configuration file.");
                return;
            } 

            if (!File.Exists(_inputFile))
            {
                Console.WriteLine($"Input file ({_inputFile}) does not exist.");
                return;
            }

            ParseInputFile();
        }

        static bool ParseConfiguration()
        {
            try
            {
                Dictionary<string, object> configuration;
                using (var sr = new StreamReader(new FileStream(CONFIG_FILE, FileMode.Open))) 
                {
                    configuration = JsonSerializer.Deserialize<Dictionary<string, object>>(sr.ReadToEnd());
                }
                _inputFile = configuration["inputFile"].ToString();
            }
            catch (Exception e)
            {
                Console.WriteLine($"Error parsing configuration file: {e.Message}");
                return false;
            }
            return true;
        }

        static bool ParseInputFile()
        {
            var csvParser = new CsvParser<Selection>(new CsvParserOptions(true, ','), new SelectionMapping());
            var result = csvParser
                .ReadFromFile(_inputFile, Encoding.ASCII)
                .ToList();

            return true;
        }

    }
}
