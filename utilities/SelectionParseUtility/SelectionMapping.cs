using System;

using TinyCsvParser;
using TinyCsvParser.Mapping;


namespace SelectionParseUtility
{
    internal class SelectionMapping : CsvMapping<Selection>
    {
        public SelectionMapping() : base()
        {
            MapProperty(0, x => x.TeamName);
            MapProperty(1, x => x.MapSelection);
            MapProperty(2, x => x.MapSelection);
            MapProperty(3, x => x.MapSelection);
            MapProperty(4, x => x.MapSelection);
        }
    }
}