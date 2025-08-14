import { Editor } from "@tinymce/tinymce-react";

import "tinymce/tinymce";
import "tinymce/icons/default";
import "tinymce/themes/silver";
import "tinymce/models/dom";
import "tinymce/plugins/table";
import "tinymce/plugins/lists";
import "tinymce/plugins/code";
import "tinymce/plugins/autoresize";
import "tinymce/plugins/importcss";

interface Props {
  value: string;
}

const TinyEditor = ({ value }: Props) => {
  return (
    <Editor
      licenseKey="gpl"
      initialValue={value}
      init={{
        height: 360,
        menubar: "file edit view insert format tools",
        plugins: "table lists code autoresize importcss",
        toolbar:
          "undo redo | styles | bold italic underline | bullist numlist | table | removeformat | code",
        skin_url: "/tinymce/skins/ui/oxide",
        content_css: [
          "/tinymce/skins/ui/oxide/content.min.css",
          "/tinymce/skins/content/default/content.min.css",
        ],

        valid_elements: "*[*]",
        extended_valid_elements: "*[*]",
      }}
    />
  );
};

export default TinyEditor;
