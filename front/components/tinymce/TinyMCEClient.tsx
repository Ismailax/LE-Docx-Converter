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

const basePath = process.env.NEXT_PUBLIC_APP_BASEPATH || "";
const withBase = (p: string) => `${basePath}${p.startsWith("/") ? p : `/${p}`}`;

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
        skin_url: withBase("/tinymce/skins/ui/oxide"),
        content_css: [
          withBase("/tinymce/skins/ui/oxide/content.min.css"),
          withBase("/tinymce/skins/content/default/content.min.css"),
        ],

        valid_elements: "*[*]",
        extended_valid_elements: "*[*]",
      }}
    />
  );
};

export default TinyEditor;
