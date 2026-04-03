-- Lua filter: map :::ClassName fenced divs to Word custom-style
-- Example: :::Warning → custom-style="Warning"

function Div(el)
  if #el.classes > 0 and not el.attributes['custom-style'] then
    el.attributes['custom-style'] = el.classes[1]
    return el
  end
end
