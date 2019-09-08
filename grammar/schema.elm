import Dict exposing (Dict)

type Json = JsonNumber Float
    | JsonString String
    | JsonBool Bool
    | JsonList (List Json)
    | JsonMap (Dict String Json)

type alias Tag = String

type Attrs = TextContent String
    | Dict String Json

type Elem = FullElem Tag Attrs (List Elem)
    | NoChildren Tag Attrs
    | NoAttrs Tag (List Elem)

type alias Selector = String
type alias Styles = Dict String String
type alias ValueQuery = List String
type alias OnWhat = String
type alias Where = Int

type Command = Subscribe Selector OnWhat ValueQuery
    | PostElem Selector Where Elem
    | PutElem Selector Elem
    | DeleteElem Selector
    | PatchAttrs Selector Attrs
    | PatchStyles Selector Styles
    | Close


type alias ValueAnswer = Dict String Json

type Event = Hi | Bye
    | UI Selector OnWhat ValueAnswer
    | Err String
