#include "chatwindow.h"
#include "ui_chatwindow.h"

ChatWindow::ChatWindow(QString eventID, QString uid, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::ChatWindow), eventID(eventID), uid(uid), cli(settings.value("host").toString(), settings.value("port").toString())
{
    ui->setupUi(this);
    loadMessages();
}

ChatWindow::~ChatWindow()
{
    delete ui;
}

void ChatWindow::on_sendButton_clicked()
{
    ui->sendButton->setDisabled(true);

    bool ok = cli.sendMessage(eventID, ui->msgEdit->toPlainText(), uid);
    if (!ok) {
        QMessageBox::warning(this, "Ошибка", "Ошибка при отправке сообщения: проверьте связь с сервером");
    }
    ui->sendButton->setEnabled(true);
    loadMessages();
    ui->msgEdit->clear();
}

void ChatWindow::loadMessages() {
    {
        QWidget* prev = ui->messages->widget();
        if (prev != nullptr) delete prev;
    }

    QVector<Message> chat = cli.getMessages(eventID);
    if (chat.isEmpty()) {
        qDebug() << "empty chat loaded";
        return;
    }
    QWidget *central = new QWidget;
    QVBoxLayout* layout = new QVBoxLayout(central);
    for (const Message& msg : chat) {
        message_entry* msgview = new message_entry;
        msgview->setAttributes(msg.sender, msg.content);
        layout->addWidget(msgview);
    }
    central->setLayout(layout);
    ui->messages->setWidget(central);
}
