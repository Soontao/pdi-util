//
// Copyright 2020 FoxyUtils ehf. All rights reserved.
//
// This is a commercial product and requires a license to operate.
// A trial license can be obtained at https://unidoc.io
//
// DO NOT EDIT: generated by unitwist Go source code obfuscator.
//
// Use of this source code is governed by the UniDoc End User License Agreement
// terms that can be accessed at https://unidoc.io/eula/

package zippkg ;import (_ab "archive/zip";_gf "bytes";_c "encoding/xml";_bd "fmt";_cb "github.com/unidoc/unioffice";_ca "github.com/unidoc/unioffice/algo";_ee "github.com/unidoc/unioffice/common/tempstorage";_gcd "github.com/unidoc/unioffice/schema/soo/pkg/relationships";_gc "io";_ba "path";_b "sort";_e "strings";_g "time";);

// OnNewRelationshipFunc is called when a new relationship has been discovered.
//
// target is a resolved path that takes into account the location of the
// relationships file source and should be the path in the zip file.
//
// files are passed so non-XML files that can't be handled by AddTarget can be
// decoded directly (e.g. images)
//
// rel is the actual relationship so its target can be modified if the source
// target doesn't match where unioffice will write the file (e.g. read in
// 'xl/worksheets/MyWorksheet.xml' and we'll write out
// 'xl/worksheets/sheet1.xml')
type OnNewRelationshipFunc func (_f *DecodeMap ,_bc ,_baf string ,_ae []*_ab .File ,_gcg *_gcd .Relationship ,_bg Target )error ;type Target struct{Path string ;Typ string ;Ifc interface{};Index uint32 ;};

// AddFileFromBytes takes a byte array and adds it at a given path to a zip file.
func AddFileFromBytes (z *_ab .Writer ,zipPath string ,data []byte )error {_da ,_cca :=z .Create (zipPath );if _cca !=nil {return _bd .Errorf ("e\u0072\u0072\u006f\u0072 c\u0072e\u0061\u0074\u0069\u006e\u0067 \u0025\u0073\u003a\u0020\u0025\u0073",zipPath ,_cca );};_ ,_cca =_gc .Copy (_da ,_gf .NewReader (data ));return _cca ;};

// SelfClosingWriter wraps a writer and replaces XML tags of the
// type <foo></foo> with <foo/>
type SelfClosingWriter struct{W _gc .Writer ;};

// AddFileFromDisk reads a file from internal storage and adds it at a given path to a zip file.
// TODO: Rename to AddFileFromStorage in next major version release (v2).
// NOTE: If disk storage cannot be used, memory storage can be used instead by calling memstore.SetAsStorage().
func AddFileFromDisk (z *_ab .Writer ,zipPath ,storagePath string )error {_cdfg ,_bbd :=z .Create (zipPath );if _bbd !=nil {return _bd .Errorf ("e\u0072\u0072\u006f\u0072 c\u0072e\u0061\u0074\u0069\u006e\u0067 \u0025\u0073\u003a\u0020\u0025\u0073",zipPath ,_bbd );};_eb ,_bbd :=_ee .Open (storagePath );if _bbd !=nil {return _bd .Errorf ("e\u0072r\u006f\u0072\u0020\u006f\u0070\u0065\u006e\u0069n\u0067\u0020\u0025\u0073: \u0025\u0073",storagePath ,_bbd );};defer _eb .Close ();_ ,_bbd =_gc .Copy (_cdfg ,_eb );return _bbd ;};

// Decode unmarshals the content of a *zip.File as XML to a given destination.
func Decode (f *_ab .File ,dest interface{})error {_cg ,_abe :=f .Open ();if _abe !=nil {return _bd .Errorf ("e\u0072r\u006f\u0072\u0020\u0072\u0065\u0061\u0064\u0069n\u0067\u0020\u0025\u0073: \u0025\u0073",f .Name ,_abe );};defer _cg .Close ();_gb :=_c .NewDecoder (_cg );if _aba :=_gb .Decode (dest );_aba !=nil {return _bd .Errorf ("e\u0072\u0072\u006f\u0072 d\u0065c\u006f\u0064\u0069\u006e\u0067 \u0025\u0073\u003a\u0020\u0025\u0073",f .Name ,_aba );};if _ga ,_ge :=dest .(*_gcd .Relationships );_ge {for _eed ,_ddc :=range _ga .Relationship {switch _ddc .TypeAttr {case _cb .OfficeDocumentTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .OfficeDocumentType ;case _cb .StylesTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .StylesType ;case _cb .ThemeTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .ThemeType ;case _cb .ControlTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .ControlType ;case _cb .SettingsTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .SettingsType ;case _cb .ImageTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .ImageType ;case _cb .CommentsTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .CommentsType ;case _cb .ThumbnailTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .ThumbnailType ;case _cb .DrawingTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .DrawingType ;case _cb .ChartTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .ChartType ;case _cb .ExtendedPropertiesTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .ExtendedPropertiesType ;case _cb .CustomXMLTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .CustomXMLType ;case _cb .WorksheetTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .WorksheetType ;case _cb .SharedStringsTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .SharedStringsType ;case _cb .TableTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .TableType ;case _cb .HeaderTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .HeaderType ;case _cb .FooterTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .FooterType ;case _cb .NumberingTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .NumberingType ;case _cb .FontTableTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .FontTableType ;case _cb .WebSettingsTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .WebSettingsType ;case _cb .FootNotesTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .FootNotesType ;case _cb .EndNotesTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .EndNotesType ;case _cb .SlideTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .SlideType ;case _cb .VMLDrawingTypeStrict :_ga .Relationship [_eed ].TypeAttr =_cb .VMLDrawingType ;};};_b .Slice (_ga .Relationship ,func (_fba ,_efe int )bool {_dge :=_ga .Relationship [_fba ];_cac :=_ga .Relationship [_efe ];return _ca .NaturalLess (_dge .IdAttr ,_cac .IdAttr );});};return nil ;};

// RelationsPathFor returns the relations path for a given filename.
func RelationsPathFor (path string )string {_dg :=_e .Split (path ,"\u002f");_cfb :=_e .Join (_dg [0:len (_dg )-1],"\u002f");_ef :=_dg [len (_dg )-1];_cfb +="\u002f_\u0072\u0065\u006c\u0073\u002f";_ef +="\u002e\u0072\u0065l\u0073";return _cfb +_ef ;};const XMLHeader ="\u003c\u003f\u0078\u006d\u006c\u0020\u0076e\u0072\u0073\u0069o\u006e\u003d\u00221\u002e\u0030\"\u0020\u0065\u006e\u0063\u006f\u0064i\u006eg=\u0022\u0055\u0054\u0046\u002d\u0038\u0022\u0020\u0073\u0074\u0061\u006e\u0064\u0061\u006c\u006f\u006e\u0065\u003d\u0022\u0079\u0065\u0073\u0022\u003f\u003e"+"\u000a";var _bec =[]byte {'\r','\n'};func (_fbad SelfClosingWriter )Write (b []byte )(int ,error ){_gbd :=0;_eec :=0;for _bab :=0;_bab < len (b )-2;_bab ++{if b [_bab ]=='>'&&b [_bab +1]=='<'&&b [_bab +2]=='/'{_fgc :=[]byte {};_fbb :=_bab ;for _gd :=_bab ;_gd >=0;_gd --{if b [_gd ]==' '{_fbb =_gd ;}else if b [_gd ]=='<'{_fgc =b [_gd +1:_fbb ];break ;};};_bbdg :=[]byte {};for _abd :=_bab +3;_abd < len (b );_abd ++{if b [_abd ]=='>'{_bbdg =b [_bab +3:_abd ];break ;};};if !_gf .Equal (_fgc ,_bbdg ){continue ;};_cbb ,_cbbf :=_fbad .W .Write (b [_gbd :_bab ]);if _cbbf !=nil {return _eec +_cbb ,_cbbf ;};_eec +=_cbb ;_ ,_cbbf =_fbad .W .Write (_abea );if _cbbf !=nil {return _eec ,_cbbf ;};_eec +=3;for _gdb :=_bab +2;_gdb < len (b )&&b [_gdb ]!='>';_gdb ++{_eec ++;_gbd =_gdb +2;_bab =_gbd ;};};};_dgef ,_bdb :=_fbad .W .Write (b [_gbd :]);return _dgef +_eec ,_bdb ;};func (_cag *DecodeMap )RecordIndex (path string ,idx int ){_cag ._d [path ]=idx };

// Decode loops decoding targets registered with AddTarget and calling th
func (_abf *DecodeMap )Decode (files []*_ab .File )error {_bbf :=1;for _bbf > 0{for len (_abf ._cf )> 0{_fdc :=_abf ._cf [len (_abf ._cf )-1];_abf ._cf =_abf ._cf [0:len (_abf ._cf )-1];_fb :=_fdc .Ifc .(*_gcd .Relationships );for _ ,_dd :=range _fb .Relationship {_cdf ,_ :=_abf ._aa [_fb ];_abf ._be (_abf ,_cdf +_dd .TargetAttr ,_dd .TypeAttr ,files ,_dd ,_fdc );};};for _db ,_ec :=range files {if _ec ==nil {continue ;};if _af ,_ace :=_abf ._fg [_ec .Name ];_ace {delete (_abf ._fg ,_ec .Name );if _eg :=Decode (_ec ,_af .Ifc );_eg !=nil {return _eg ;};files [_db ]=nil ;if _cc ,_cda :=_af .Ifc .(*_gcd .Relationships );_cda {_abf ._cf =append (_abf ._cf ,_af );_cfe ,_ :=_ba .Split (_ba .Clean (_ec .Name +"\u002f\u002e\u002e\u002f"));_abf ._aa [_cc ]=_cfe ;_bbf ++;};};};_bbf --;};return nil ;};

// ExtractToDiskTmp extracts a zip file to a temporary file in a given path,
// returning the name of the file.
func ExtractToDiskTmp (f *_ab .File ,path string )(string ,error ){_ea ,_bdc :=_ee .TempFile (path ,"\u007a\u007a");if _bdc !=nil {return "",_bdc ;};defer _ea .Close ();_fdg ,_bdc :=f .Open ();if _bdc !=nil {return "",_bdc ;};defer _fdg .Close ();_ ,_bdc =_gc .Copy (_ea ,_fdg );if _bdc !=nil {return "",_bdc ;};return _ea .Name (),nil ;};

// AddTarget allows documents to register decode targets. Path is a path that
// will be found in the zip file and ifc is an XML element that the file will be
// unmarshaled to.  filePath is the absolute path to the target, ifc is the
// object to decode into, sourceFileType is the type of file that the reference
// was discovered in, and index is the index of the source file type.
func (_aef *DecodeMap )AddTarget (filePath string ,ifc interface{},sourceFileType string ,idx uint32 )bool {if _aef ._fg ==nil {_aef ._fg =make (map[string ]Target );_aef ._aa =make (map[*_gcd .Relationships ]string );_aef ._caf =make (map[string ]struct{});_aef ._d =make (map[string ]int );};_bb :=_ba .Clean (filePath );if _ ,_fd :=_aef ._caf [_bb ];_fd {return false ;};_aef ._caf [_bb ]=struct{}{};_aef ._fg [_bb ]=Target {Path :filePath ,Typ :sourceFileType ,Ifc :ifc ,Index :idx };return true ;};

// DecodeMap is used to walk a tree of relationships, decoding files and passing
// control back to the document.
type DecodeMap struct{_fg map[string ]Target ;_aa map[*_gcd .Relationships ]string ;_cf []Target ;_be OnNewRelationshipFunc ;_caf map[string ]struct{};_d map[string ]int ;};func (_cd *DecodeMap )IndexFor (path string )int {return _cd ._d [path ]};

// MarshalXML creates a file inside of a zip and marshals an object as xml, prefixing it
// with a standard XML header.
func MarshalXML (z *_ab .Writer ,filename string ,v interface{})error {_fge :=&_ab .FileHeader {};_fge .Method =_ab .Deflate ;_fge .Name =filename ;_fge .SetModTime (_g .Now ());_bf ,_fe :=z .CreateHeader (_fge );if _fe !=nil {return _bd .Errorf ("\u0063\u0072\u0065\u0061ti\u006e\u0067\u0020\u0025\u0073\u0020\u0069\u006e\u0020\u007a\u0069\u0070\u003a\u0020%\u0073",filename ,_fe );};_ ,_fe =_bf .Write ([]byte (XMLHeader ));if _fe !=nil {return _bd .Errorf ("\u0063\u0072e\u0061\u0074\u0069\u006e\u0067\u0020\u0078\u006d\u006c\u0020\u0068\u0065\u0061\u0064\u0065\u0072\u0020\u0074\u006f\u0020\u0025\u0073: \u0025\u0073",filename ,_fe );};if _fe =_c .NewEncoder (SelfClosingWriter {_bf }).Encode (v );_fe !=nil {return _bd .Errorf ("\u006d\u0061\u0072\u0073\u0068\u0061\u006c\u0069\u006e\u0067\u0020\u0025s\u003a\u0020\u0025\u0073",filename ,_fe );};_ ,_fe =_bf .Write (_bec );return _fe ;};var _abea =[]byte {'/','>'};func MarshalXMLByType (z *_ab .Writer ,dt _cb .DocType ,typ string ,v interface{})error {_gfeb :=_cb .AbsoluteFilename (dt ,typ ,0);return MarshalXML (z ,_gfeb ,v );};func MarshalXMLByTypeIndex (z *_ab .Writer ,dt _cb .DocType ,typ string ,idx int ,v interface{})error {_ddg :=_cb .AbsoluteFilename (dt ,typ ,idx );return MarshalXML (z ,_ddg ,v );};

// SetOnNewRelationshipFunc sets the function to be called when a new
// relationship has been discovered.
func (_fgf *DecodeMap )SetOnNewRelationshipFunc (fn OnNewRelationshipFunc ){_fgf ._be =fn };