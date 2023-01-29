using System;
using System.Collections.Generic;

namespace SelectionParseUtility
{
    internal class Selection 
    {  
        public string TeamName { get; set; }
        public string MapSelection 
        { 
            get
            {
                return string.Empty;
            }
            set
            {
                Selections.Add(value);
            }
        } // only used for mapping CSV
        public List<string> Selections { get; set; } = new List<string>();
    }
}