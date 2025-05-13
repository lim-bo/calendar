#ifndef MESSAGE_ENTRY_H
#define MESSAGE_ENTRY_H

#include <QWidget>

namespace Ui {
class message_entry;
}

class message_entry : public QWidget
{
    Q_OBJECT

public:
    explicit message_entry(QWidget *parent = nullptr);
    ~message_entry();
    void setAttributes(QString sender, QString content);
private:
    Ui::message_entry *ui;
};

#endif // MESSAGE_ENTRY_H
